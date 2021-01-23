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

        public Base(DF.Game.API.Data.Curve pb)
        {
            ID = new DF.Game.ID.EntityID(pb.EntityId);
            Type = pb.Type;
            Property = pb.Property;
            Tick = new DF.Game.ID.Tick(pb.Tick);
        }

        public DF.Game.ID.EntityID ID
        {
            get => _entityID;
            private set { _entityID = value; }
        }
        public DF.Game.API.Constants.CurveType Type
        {
            get => _type;
            private set { _type = value; }
        }
        public DF.Game.API.Constants.EntityProperty Property
        {
            get => _property;
            private set { _property = value; }
        }
        public DF.Game.ID.Tick Tick
        {
            get => _tick;
            protected set { _tick = value; }
        }
    }
}
