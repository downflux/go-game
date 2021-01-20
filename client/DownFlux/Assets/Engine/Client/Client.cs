namespace DF.Game.Client
{
    public class Client
    {
        private DF.Game.API.API.DownFlux.DownFluxClient _c;
        private System.Threading.CancellationTokenSource _ctSource;
        private System.Threading.CancellationToken _ct;
        private string _cid;

        public Client(string server, string id)
        {
            _c = new DF.Game.API.API.DownFlux.DownFluxClient(
                new Grpc.Core.Channel(server, Grpc.Core.ChannelCredentials.Insecure));
            _ctSource = new System.Threading.CancellationTokenSource();
            _ct = _ctSource.Token;
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
    }
}
