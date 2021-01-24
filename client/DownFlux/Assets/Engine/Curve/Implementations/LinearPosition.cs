namespace DF.Game.Curve.LinearPosition
{
    public class LinearPosition : DF.Game.Curve.Base
    {
        public LinearPosition(DF.Game.API.Data.Curve pb) : base(pb)
        {

        }

        public void Add(DF.Game.ID.Tick tick, object data)
        {
            var p = data as DF.Game.API.Data.Position;
        }

        public object Get(DF.Game.ID.Tick tick)
        {
            return new System.Exception();
        }

        public void Merge(LinearPosition curve)
        {
            if (Type != curve.Type)
            {
                throw new DF.Game.Exception.MergeException(
                    System.String.Format("Cannot merge curves: {0} != {0}.", Type, curve.Type));
            }
        }
    }
}