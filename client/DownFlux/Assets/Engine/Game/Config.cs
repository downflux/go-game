namespace DF.Game
{
    public class Config
    {
        private System.TimeSpan _serverBootSleepDuration;
        private System.TimeSpan _entityListAcquireDuration;
        
        public Config(
            System.TimeSpan serverBootSleepDuration,
            System.TimeSpan entityListAcquireDuration
        )
        {
            ServerBootSleepDuration = serverBootSleepDuration;
            EntityListAcquireDuration = entityListAcquireDuration;
        }

        public System.TimeSpan ServerBootSleepDuration
        {
            get => _serverBootSleepDuration;
            private set { _serverBootSleepDuration = value; }
        }

        public System.TimeSpan EntityListAcquireDuration
        {
            get => _entityListAcquireDuration;
            private set { _entityListAcquireDuration = value; }
        }
    }
}