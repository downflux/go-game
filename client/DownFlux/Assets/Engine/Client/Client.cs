using System.Linq;

namespace DF.Game.Client
{
    public class Client
    {
        private System.TimeSpan _sleepDuration;
        private DF.Game.API.API.DownFlux.DownFluxClient _c;
        private System.Threading.CancellationTokenSource _ctSource;
        private System.Threading.CancellationToken _ct;
        private string _cid;

        public Client(string server, System.TimeSpan sleepDuration)
        {
            _c = new DF.Game.API.API.DownFlux.DownFluxClient(
                new Grpc.Core.Channel(server, Grpc.Core.ChannelCredentials.Insecure));
            _ctSource = new System.Threading.CancellationTokenSource();
            _ct = _ctSource.Token;
            _sleepDuration = sleepDuration;
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

        public DF.Game.ServerStatus.ServerStatus WaitForBoot()
        {
            var status = new DF.Game.ServerStatus.ServerStatus(new DF.Game.API.Data.ServerStatus());
            while (!status.IsStarted)
            {
                System.Threading.Thread.Sleep(_sleepDuration);
                status = new DF.Game.ServerStatus.ServerStatus(
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
            var req = new DF.Game.API.API.MoveRequest{
                ClientId = ID,
                Tick = tick.Double,
                EntityIds = { eids },
                Destination = destination,
                MoveType = moveType,
            };
            return _c.Move(req);
        }
    }
}