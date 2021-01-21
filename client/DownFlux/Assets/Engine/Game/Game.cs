namespace DF.Game {
    public class Game {
        private DF.Game.State.State _state;
        private DF.Game.ServerStatus.ServerStatus _status;
        private DF.Game.Client.Client _client;
        private DF.Game.Config _config;

        public Game(string server) {
            _client = new DF.Game.Client.Client(server, _config.ServerBootSleepDuration);
            _client.Connect();
            _client.WaitForBoot();
        }
    }
}