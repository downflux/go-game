namespace DF.Game.Curve
{
    public class CurveOneOf : OneOf.OneOfBase<
        DF.Game.Curve.LinearPosition.LinearPosition
    >
    {
        public static CurveOneOf Import(DF.Game.API.Data.Curve c)
        {
            switch (c.Property)
            {
                case DF.Game.API.Constants.EntityProperty.Position:
                    return new CurveOneOf(
                        new DF.Game.Curve.LinearPosition.LinearPosition(c)
                    );
                default:
                    throw new System.NotImplementedException(
                        System.String.Format("Curve type {0} does not have a corresponding client definition.", c.Property));
            }
        }
        public CurveOneOf(OneOf.OneOf<
            DF.Game.Curve.LinearPosition.LinearPosition
        > c) : base(c) { }

        public DF.Game.API.Constants.EntityProperty Property
        {
            get
            {
                return Match(
                    linearPosition => DF.Game.API.Constants.EntityProperty.Position
                );
            }
        }

        public (bool ok, DF.Game.Curve.LinearPosition.LinearPosition) TryGetLinearPosition()
        {
            return Match(
                linearPosition => (true, linearPosition)
            );
        }

        public void Merge(CurveOneOf other)
        {
            Switch(
                linearPosition =>
                {
                    var (ok, c) = other.TryGetLinearPosition();
                    if (ok)
                    {
                        linearPosition.Merge(c);
                    }
                }
            );
        }
    }
}