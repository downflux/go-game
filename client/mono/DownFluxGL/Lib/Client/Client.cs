namespace DF {
  namespace Client {
    // double here represents the server tick.
    // TODO(minkezhang): Remove the tuple -- this can be a naked queue.
    using CurveStream = System.Collections.Generic.Queue<
      (double, DF.Curve.Curve)>;
    using EntityStream = System.Collections.Generic.Queue<
      (double, DF.Entity.Entity)>;

    public class Client {
      private DF.Game.API.API.DownFlux.DownFluxClient _client;
      private string _id;
      private System.Threading.CancellationTokenSource _ctSource;
      private System.Threading.CancellationToken _ct;

      public void Cancel() { _ctSource.Cancel(); }

      // TODO(minkezhang): Add entity list.

      private System.Threading.ReaderWriterLock _curvesMutex;
      private CurveStream _curves;

      private System.Threading.ReaderWriterLock _entitiesMutex;
      private EntityStream _entities;

      public string ID { get => _id; }

      public Client(Grpc.Core.Channel channel) {
        _client = new DF.Game.API.API.DownFlux.DownFluxClient(channel);
        _ctSource = new System.Threading.CancellationTokenSource();
        _ct = _ctSource.Token;

        _curvesMutex = new System.Threading.ReaderWriterLock();
        _curves = new CurveStream();

        _entitiesMutex = new System.Threading.ReaderWriterLock();
        _entities = new EntityStream();
      }

      public string Connect() {
        var resp = _client.AddClient(new DF.Game.API.API.AddClientRequest());
        _id = resp.ClientId;
        return ID;
      }

      // Returns the current buffer of streamed messages.
      // TODO(minkezhang): Rename Curves or something.
      public CurveStream Data {
        get {
          _curvesMutex.AcquireReaderLock(1000);  // 1 sec
          try {
            var ret = _curves;
            _curves = new CurveStream();
            return ret;
          } finally {
            _curvesMutex.ReleaseReaderLock();
          }
        }
      }

      public EntityStream Entities {
        get {
          _entitiesMutex.AcquireReaderLock(1000);
          try {
            var ret = _entities;
            _entities = new EntityStream();
            return ret;
          } finally {
            _entitiesMutex.ReleaseReaderLock();
          }
        }
      }

      public DF.Game.API.Data.ServerStatus GetStatus() {
        return _client.GetStatus(new DF.Game.API.API.GetStatusRequest{}).Status;
      }

      public void Move(
        double tick,
        System.Collections.Generic.List<string> entityIDs,
        DF.Game.API.Data.Position destination,
        DF.Game.API.Constants.MoveType moveType) {
        _client.Move(new DF.Game.API.API.MoveRequest{
          ClientId = ID,
          Tick = tick,
          EntityIds = { entityIDs },
          Destination = destination,
          MoveType = moveType,
        });
      }

      public async System.Threading.Tasks.Task StreamCurvesLoop(double tick) {
        using (var call = _client.StreamCurves(new DF.Game.API.API.StreamCurvesRequest{
          ClientId = ID,
          Tick = tick,
        })) {
          var s = call.ResponseStream;
          try {
            while (await s.MoveNext(_ct)) {
              System.Console.Error.WriteLine("StreamCurvesLoop: RECEIVED");
              var resp = s.Current;
              System.Console.Error.WriteLine(resp);

              _curvesMutex.AcquireWriterLock(1000);  // 1 sec
              try {
                foreach (var curvePB in resp.Curves) {
                  try {
                    _curves.Enqueue((resp.Tick, DF.Curve.Curve.Import(curvePB)));
                  } catch (System.ArgumentException) {
                     // TODO(minkezhang): Log this to some file.
                  }
                }
              } finally {
                _curvesMutex.ReleaseWriterLock();
              }
              // TODO(minkezhang): Parallelize this.
              _entitiesMutex.AcquireWriterLock(1000);  // 1 sec
              try {
                foreach (var entityPB in resp.Entities) {
                  try {
                    _entities.Enqueue((resp.Tick, DF.Entity.Entity.Import(entityPB)));
                  } catch (System.ArgumentException) {
                  }
                }
              } finally {
                _entitiesMutex.ReleaseWriterLock();
              }
            }
          } catch (System.Threading.Tasks.TaskCanceledException) {
          }
        }
      }
    }
  }
}
