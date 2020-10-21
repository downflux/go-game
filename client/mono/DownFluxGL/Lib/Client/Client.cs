namespace DownFluxClient {
    public class Client {
        private DF.Game.API.API.DownFlux.DownFluxClient _client;

        public Client(Grpc.Core.Channel channel) {
            _client = new DF.Game.API.API.DownFlux.DownFluxClient(channel);
        }
    }
}
