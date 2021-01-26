namespace DF.Game
{
    public class Game
    {
        private DF.Game.Entity.List _entities;
        private DF.Game.Status.Status _status;
        private DF.Game.Client.Client _client;
        private DF.Game.Config _config;
        private DF.Game.Entity.Listener.LastState _listener;

        public Game(
            string server, DF.Game.Config config,
            System.Threading.CancellationToken ct,
            System.Action<DF.Game.Entity.Entity> cb)
        {
            _entities = new DF.Game.Entity.List(cb);

            _config = config;
            _listener = new DF.Game.Entity.Listener.LastState(config.ListenerAcquireTimeout);

            _client = new DF.Game.Client.Client(server, _config.ServerBootSleepDuration);
            _client.Connect();
            _status = _client.WaitForBoot(_config.ServerBootSleepDuration);

            // Initializes StreamData thread, which constantly updates the
            // _listener while not mutating the actual game state.
            _client.StreamData(ct, _status.Tick, _listener.Listen);
        }

        public DF.Game.ID.Tick Tick
        {
            get => _status.Tick;
        }

        // Merge merges the internal state cache with any potential updates
        // from the gRPC server.
        //
        // This is expensive and should only be called once every several
        // ticks.
        public void Merge()
        {
            var s = _listener.Pop();
            if (s == null)
            {
                return;
            }
            _entities.Merge(s);
        }
    }
}