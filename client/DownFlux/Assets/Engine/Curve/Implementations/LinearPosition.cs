namespace DF.Game.Curve.LinearPosition
{
    public class LinearPosition : DF.Game.Curve.Base
    {
        private DF.Game.Data.Data<DF.Game.API.Data.Position> _data;

        public LinearPosition(DF.Game.API.Data.Curve pb) : base(pb)
        {
            var data = new System.Collections.Generic.List<DF.Game.Data.Datum<DF.Game.API.Data.Position>>();
            if (pb.Property != DF.Game.API.Constants.EntityProperty.Position)
            {
                throw new System.ArgumentException(
                    System.String.Format(
                        "Cannot construct a LinearPosition curve with property type {0} != {1}",
                        pb.Property,
                        DF.Game.API.Constants.EntityProperty.Position));
            }
            foreach (var d in pb.Data)
            {
                // Assume curve data is in tick order; if we cannot assume
                // this condition we need to manually sort the list before
                // insertion.
                data.Add(
                    new DF.Game.Data.Datum<DF.Game.API.Data.Position>(
                        new DF.Game.ID.Tick(d.Tick),
                        d.PositionDatum));
            }

            _data = new DF.Game.Data.Data<DF.Game.API.Data.Position>(data);
        }

        public DF.Game.API.Data.Position Get(DF.Game.ID.Tick tick)
        {
            throw new System.NotImplementedException();
        }

        public void Merge(LinearPosition curve)
        {
            if (Type != curve.Type)
            {
                throw new System.ArgumentException(
                    System.String.Format("Cannot merge curves: {0} != {0}.", Type, curve.Type));
            }
        }
    }
}