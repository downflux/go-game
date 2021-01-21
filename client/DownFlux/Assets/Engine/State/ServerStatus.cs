namespace DF.Game.ServerStatus
{
    public class ServerStatus
    {
        private DF.Game.API.Data.ServerStatus _pb;
        public System.TimeSpan TickDuration { get => _pb.TickDuration.ToTimeSpan(); }
        public DF.Game.ID.Tick Tick { get => new DF.Game.ID.Tick(_pb.Tick); }
        public bool IsStarted { get => _pb.IsStarted; }
        public System.DateTime StartTime { get => _pb.StartTime.ToDateTime(); }

        public ServerStatus(DF.Game.API.Data.ServerStatus pb) {
            _pb = pb;
        }
    }
}