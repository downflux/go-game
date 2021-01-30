[assembly: System.Runtime.CompilerServices.InternalsVisibleTo("EngineTestAssembly")]

namespace DF.Test.Data
{
    internal static class StatusTestData
    {
        internal static DF.Game.API.Data.ServerStatus pb
        {
            get => new DF.Game.API.Data.ServerStatus
            {
                Tick = 10,
                StartTime = new Google.Protobuf.WellKnownTypes.Timestamp
                {
                    Seconds = 946684800  // Jan 1, 2000
                },
                TickDuration = new Google.Protobuf.WellKnownTypes.Duration
                {
                    Nanos = 100000000  // 100ms
                },
                IsStarted = true
            };
        }
    }
}