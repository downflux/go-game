namespace DF.Game.Entity
{
    public interface IEntity
    {
        DF.Game.API.Constants.EntityType Type { get; }
        DF.Game.ID.EntityID ID { get; }
        DF.Game.Curve.List Curves { get; }
    }

    public class Entity : IEntity
    {
        DF.Game.API.Constants.EntityType _type;
        DF.Game.ID.EntityID _id;

        public Entity(DF.Game.API.Data.Entity pb) : this(
            type: pb.Type,
            id: new DF.Game.ID.EntityID(pb.EntityId)
        )
        {
        }

        public Entity(DF.Game.API.Constants.EntityType type, DF.Game.ID.EntityID id)
        {
            Type = type;
            ID = id;
        }

        public DF.Game.API.Constants.EntityType Type
        {
            get => _type;
            private set { _type = value; }
        }

        public DF.Game.ID.EntityID ID
        {
            get => _id;
            private set { _id = value; }
        }

        public DF.Game.Curve.List Curves { get; }
    }
}