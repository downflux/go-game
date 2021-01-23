namespace DF.Game.Status
{
    // Status is the current gRPC server status.
    public class Status
    {
        private DF.Game.API.Data.ServerStatus _pb;
        public System.TimeSpan TickDuration { get => _pb.TickDuration.ToTimeSpan(); }
        public DF.Game.ID.Tick Tick
        {
            get => new DF.Game.ID.Tick(
                _pb.Tick + System.DateTime.UtcNow.Subtract(StartTime).TotalMilliseconds / TickDuration.TotalMilliseconds
            );
        }
        public bool IsStarted { get => _pb.IsStarted; }
        public System.DateTime StartTime { get => _pb.StartTime.ToDateTime(); }

        public Status(DF.Game.API.Data.ServerStatus pb)
        {
            _pb = pb;
        }
    }
}