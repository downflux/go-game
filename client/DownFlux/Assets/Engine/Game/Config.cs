namespace DF.Game
{
    public class Config
    {
        private System.TimeSpan _serverBootSleepDuration;
        private System.TimeSpan _listenerAquireTimeout;
        private DF.Game.ID.Tick _updateTickDelay;

        public Config(
            System.TimeSpan serverBootSleepDuration,
            System.TimeSpan listenerAcquireTimeout,
            DF.Game.ID.Tick updateTickDelay
        )
        {
            ServerBootSleepDuration = serverBootSleepDuration;
            ListenerAcquireTimeout = listenerAcquireTimeout;
            UpdateTickDelay = updateTickDelay;
        }

        public System.TimeSpan ServerBootSleepDuration
        {
            get => _serverBootSleepDuration;
            private set { _serverBootSleepDuration = value; }
        }
        public System.TimeSpan ListenerAcquireTimeout {
            get => _listenerAquireTimeout;
            private set { _listenerAquireTimeout = value; }
        }
        public DF.Game.ID.Tick UpdateTickDelay
        {
            get => _updateTickDelay;
            private set { _updateTickDelay = value; }
        }
    }
}