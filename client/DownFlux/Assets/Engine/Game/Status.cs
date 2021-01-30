namespace DF.Game.Status
{
    // Status is the current gRPC server status.
    public class Status
    {
        private DF.Game.API.Data.ServerStatus _pb;
        public System.TimeSpan TickDuration { get => _pb.TickDuration.ToTimeSpan(); }

        public DF.Game.ID.Tick Tick
        {
            get => CalculateTick(System.DateTime.UtcNow);
        }

        internal DF.Game.ID.Tick CalculateTick(System.DateTime t)
        {
            return new DF.Game.ID.Tick(t.Subtract(StartTime).TotalMilliseconds / TickDuration.TotalMilliseconds);
        }

        public bool IsStarted { get => _pb.IsStarted; }
        public System.DateTime StartTime { get => _pb.StartTime.ToDateTime(); }

        public Status(DF.Game.API.Data.ServerStatus pb)
        {
            _pb = pb;
        }
    }
}