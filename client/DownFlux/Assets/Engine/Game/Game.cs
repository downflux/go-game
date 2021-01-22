namespace DF.Game {
    public class Game {
        private DF.Game.State.State _state;
        private DF.Game.ServerStatus.ServerStatus _status;
        private DF.Game.Client.Client _client;
        private DF.Game.Config _config;
        private DF.Game.Entity.Listener.LastState _listener;

        public Game(string server, DF.Game.Config config, System.Threading.CancellationToken ct) {
            _config = config;
            _listener = new DF.Game.Entity.Listener.LastState(config.EntityListAcquireDuration);

            _client = new DF.Game.Client.Client(server, _config.ServerBootSleepDuration);
            _client.Connect();
            _status = _client.WaitForBoot(_config.ServerBootSleepDuration);

            // Initializes StreamData thread, which constantly updates the
            // _listener while not mutating the actual game state.
            _client.StreamData(ct, _status.Tick, _listener.Listen);
        }
    }
}