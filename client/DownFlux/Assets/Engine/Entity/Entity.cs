namespace DF.Game.Entity
{
    public class Entity
    {
        DF.Game.API.Constants.EntityType _type;
        DF.Game.ID.EntityID _id;
        DF.Game.Curve.List _cs;

        public Entity(DF.Game.API.Data.Entity pb) : this(
            type: pb.Type,
            id: new DF.Game.ID.EntityID(pb.EntityId)
        )
        {
            _cs = new DF.Game.Curve.List();
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

        public DF.Game.Curve.List Curves { get => _cs; }

        public void Merge(Entity e)
        {
            if (ID != e.ID)
            {
                throw new DF.Game.Exception.MergeException(
                    System.String.Format("Cannot merge entities with mismatching IDs: {0} != {1}", ID, e.ID));
            }

            Curves.Merge(e.Curves);
        }
    }
}