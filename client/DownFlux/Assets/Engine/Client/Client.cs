using System.Linq;

namespace DF.Game.Client
{
    public class Client
    {
        private DF.Game.API.API.DownFlux.DownFluxClient _c;
        private string _cid;

        public Client(string server, System.TimeSpan sleepDuration)
        {
            _c = new DF.Game.API.API.DownFlux.DownFluxClient(
                new Grpc.Core.Channel(server, Grpc.Core.ChannelCredentials.Insecure));
        }

        public string ID
        {
            get => _cid;
            private set { _cid = value; }
        }

        public string Connect()
        {
            var resp = _c.AddClient(new DF.Game.API.API.AddClientRequest());
            ID = resp.ClientId;
            return ID;
        }

        public DF.Game.Status.Status WaitForBoot(System.TimeSpan s)
        {
            var status = new DF.Game.Status.Status(new DF.Game.API.Data.ServerStatus());
            while (!status.IsStarted)
            {
                System.Threading.Thread.Sleep(s);
                status = new DF.Game.Status.Status(
                    _c.GetStatus(new DF.Game.API.API.GetStatusRequest()).Status
                );
            }
            return status;
        }

        public DF.Game.API.API.MoveResponse Move(
            DF.Game.ID.Tick tick,
            System.Collections.Generic.List<DF.Game.ID.EntityID> entityIDs,
            DF.Game.API.Data.Position destination,
            DF.Game.API.Constants.MoveType moveType)
        {
            System.Collections.Generic.IEnumerable<string> eids = from eid in entityIDs select eid.String;
            var req = new DF.Game.API.API.MoveRequest
            {
                ClientId = ID,
                Tick = tick.Double,
                EntityIds = { eids },
                Destination = destination,
                MoveType = moveType,
            };
            return _c.Move(req);
        }

        public async void StreamData(
            System.Threading.CancellationToken ct,
            DF.Game.ID.Tick t,
            System.Action<DF.Game.API.API.StreamDataResponse> f)
        {
            var stream = _c.StreamData(
                new DF.Game.API.API.StreamDataRequest
                {
                    ClientId = ID,
                    Tick = t.Double
                }
            );

            try
            {
                while (await stream.ResponseStream.MoveNext(ct))
                {
                    f(stream.ResponseStream.Current);
                }
            }
            catch (System.Threading.Tasks.TaskCanceledException)
            {
            }
            catch (Grpc.Core.RpcException ex) when (ex.StatusCode == Grpc.Core.StatusCode.Cancelled)
            {
            }
            catch (Grpc.Core.RpcException ex) when (ex.StatusCode == Grpc.Core.StatusCode.Unavailable)
            {
                // TODO(minkezhang): Implement reconnection logic.
            }
        }
    }
}