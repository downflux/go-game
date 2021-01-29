[assembly: System.Runtime.CompilerServices.InternalsVisibleTo("EngineTestAssembly")]

namespace DF.Test.Data
{
    internal static class LinearPositionTestData
    {
        internal static DF.Game.API.Data.Curve pb
        {
            get
            {
                var pb = new DF.Game.API.Data.Curve
                {
                    Type = DF.Game.API.Constants.CurveType.LinearMove,
                    Property = DF.Game.API.Constants.EntityProperty.Position,
                    Tick = 0
                };
                pb.Data.Add(new DF.Game.API.Data.CurveDatum
                {
                    Tick = 0,
                    PositionDatum = new DF.Game.API.Data.Position { X = 0, Y = 0 }
                });
                pb.Data.Add(new DF.Game.API.Data.CurveDatum
                {
                    Tick = 10,
                    PositionDatum = new DF.Game.API.Data.Position { X = 1, Y = 1 }
                });

                return pb;
            }
        }

        internal static DF.Game.API.Data.Curve mergePB
        {
            get
            {
                var pb = new DF.Game.API.Data.Curve
                {
                    Type = DF.Game.API.Constants.CurveType.LinearMove,
                    Property = DF.Game.API.Constants.EntityProperty.Position,
                    Tick = 1
                };
                pb.Data.Add(new DF.Game.API.Data.CurveDatum
                {
                    Tick = 1,
                    PositionDatum = new DF.Game.API.Data.Position { X = 1, Y = 1 }
                });
                pb.Data.Add(new DF.Game.API.Data.CurveDatum
                {
                    Tick = 11,
                    PositionDatum = new DF.Game.API.Data.Position { X = 2, Y = 2 }
                });

                return pb;
            }
        }
    }
}