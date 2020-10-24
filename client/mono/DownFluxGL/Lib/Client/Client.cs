namespace DF {
  namespace Client {
    using StreamData = System.Collections.Generic.Queue<
      (string, OneOf.OneOf<DF.Curve.LinearMove>)>;

    public class Client {
      private DF.Game.API.API.DownFlux.DownFluxClient _client;
      private string _id;
      private System.Threading.CancellationTokenSource _ctSource;
      private System.Threading.CancellationToken _ct;

      public void Cancel() { _ctSource.Cancel(); }

      // TODO(minkezhang): Add entity list.

      private System.Threading.ReaderWriterLock _curvesMutex;
      private StreamData _curves;

      public string ID { get => _id; }

      public Client(Grpc.Core.Channel channel) {
        _client = new DF.Game.API.API.DownFlux.DownFluxClient(channel);
        _ctSource = new System.Threading.CancellationTokenSource();
        _ct = _ctSource.Token;
        _curvesMutex = new System.Threading.ReaderWriterLock();
        _curves = new StreamData();
      }

      public string Connect(string tickID) {
        var resp = _client.AddClient(new DF.Game.API.API.AddClientRequest());
        _id = resp.ClientId;
        return ID;
      }

      // Returns the current buffer of streamed messages.
      public StreamData Data {
        get {
          var ret = _curves;
          _curves = new StreamData();
          return ret;
        }
      }

      public void Move(
        string tickID,
        System.Collections.Generic.List<string> entityIDs,
        DF.Game.API.Data.Position destination,
        DF.Game.API.Constants.MoveType moveType) {
        _client.Move(new DF.Game.API.API.MoveRequest{
          ClientId = ID,
          TickId = tickID,
          EntityIds = { entityIDs },
          Destination = destination,
          MoveType = moveType,
        });
      }

      public async System.Threading.Tasks.Task StreamCurvesLoop(string tickID) {
        using (var call = _client.StreamCurves(new DF.Game.API.API.StreamCurvesRequest{
          ClientId = ID,
          TickId = tickID,
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
                    System.Console.Error.WriteLine("IMPORTED CURVE: ");
                    System.Console.Error.WriteLine(DF.Curve.Curve.Import(curvePB));
                    _curves.Enqueue((resp.TickId, DF.Curve.Curve.Import(curvePB)));
                  } catch (System.ArgumentException e) {
                     // TODO(minkezhang): Log this to some file.
                  }
                }
              } finally {
                _curvesMutex.ReleaseWriterLock();
              }
            }
          } catch (System.Threading.Tasks.TaskCanceledException) {
          }
        }
      }
    }
  }
}
