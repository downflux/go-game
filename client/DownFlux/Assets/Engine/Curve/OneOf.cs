namespace DF.Game.Curve
{
    public class CurveOneOf : OneOf.OneOfBase<
        DF.Game.Curve.LinearPosition.LinearPosition
    >
    {
        CurveOneOf(OneOf.OneOf<
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