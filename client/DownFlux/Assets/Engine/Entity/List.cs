namespace DF.Game.Entity
{
    public class List
    {
        private System.Collections.Generic.Dictionary<
            DF.Game.ID.EntityID,
            DF.Game.Entity.Entity> _entities;

        public List()
        {
            _entities = new System.Collections.Generic.Dictionary<DF.Game.ID.EntityID, DF.Game.Entity.Entity>();
        }

        public List(DF.Game.API.API.StreamDataResponse pb) : this()
        {
            foreach (var e in pb.State.Entities)
            {
                Append(new DF.Game.Entity.Entity(e));
            }
            foreach (var c in pb.State.Curves)
            {
                var eid = new DF.Game.ID.EntityID(c.EntityId);
                if (!_entities.ContainsKey(eid))
                {
                    Append(new DF.Game.Entity.Entity(eid));

                }
                // TODO(minkezhang): Make this import all curves.
                if (c.Property == DF.Game.API.Constants.EntityProperty.Position)
                {
                    Entity(eid).Curves.Add(DF.Game.Curve.CurveOneOf.Import(c));
                }
            }
        }

        private DF.Game.Entity.Entity Entity(DF.Game.ID.EntityID id)
        {
            return _entities[id];
        }

        internal void Append(Entity e)
        {
            if (_entities.ContainsKey(e.ID))
            {
                _entities[e.ID].Merge(e);
            }
            else
            {
                _entities[e.ID] = e;
            }
        }

        protected System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.Entity.Entity> Entities
        {
            get => new System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.Entity.Entity>(
                          new System.Collections.Generic.List<DF.Game.Entity.Entity>(_entities.Values));
        }
        public void Merge(List o)
        {
            foreach (var e in o.Entities)
            {
                Append(e);
            }
        }
    }
}