namespace DF.Game.Curve
{
    public interface ICurve
    {
        DF.Game.ID.EntityID ID { get; }
        DF.Game.API.Constants.CurveType Type { get; }
        DF.Game.API.Constants.EntityProperty Property { get; }
        DF.Game.ID.Tick Tick { get; }
        void Add(DF.Game.ID.Tick t, object v);
        object Get(DF.Game.ID.Tick tick);
        void Merge(ICurve c);

    }

    public abstract class Base
    {
        private DF.Game.ID.EntityID _entityID;
        private DF.Game.API.Constants.CurveType _type;
        private DF.Game.API.Constants.EntityProperty _property;
        private DF.Game.ID.Tick _tick;
        public DF.Game.ID.EntityID ID { get => _entityID; }
        public DF.Game.API.Constants.CurveType Type { get => _type; }
        public DF.Game.API.Constants.EntityProperty Property { get => _property; }
        public DF.Game.ID.Tick Tick { get => _tick; }
    }

    public class LinearPosition : Base, ICurve
    {
        public void Add(DF.Game.ID.Tick tick, object data)
        {
            var p = (DF.Game.API.Data.Position)data;
        }

        public object Get(DF.Game.ID.Tick tick) {
            return null;
        }

        public void Merge(ICurve curve) {
            if (Type != curve.Type) {
                throw new DF.Game.Exception.MergeException(
                    System.String.Format("Cannot merge curves: {0} != {0}.", Type, curve.Type));
            }
        }
    }
}
