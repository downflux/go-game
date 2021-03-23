using EntityCallback = System.Action<DF.Game.Entity.Entity>;
namespace DF.Game.Entity
{
    public delegate bool F(DF.Game.Entity.Entity e);

    public class List
    {
        private System.Collections.Generic.Dictionary<
            DF.Game.ID.EntityID,
            DF.Game.Entity.Entity> _entities;
        private EntityCallback _newEntityCallback;

        public List(EntityCallback callback)
        {
            _entities = new System.Collections.Generic.Dictionary<DF.Game.ID.EntityID, DF.Game.Entity.Entity>();
            _newEntityCallback = callback;
        }

        public List(
            DF.Game.API.API.StreamDataResponse pb
        ) : this(pb, delegate (DF.Game.Entity.Entity _) { })
        {
        }

        public List(
            DF.Game.API.API.StreamDataResponse pb,
            EntityCallback callback
        ) : this(callback)
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
                if (DF.Game.Curve.CurveOneOf.ImplementedProperties.Contains(c.Property))
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
                _newEntityCallback(e);
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

        public List Filter(F f)
        {
            var l = new List(null, delegate (DF.Game.Entity.Entity _) { });
            foreach (var e in _entities.Values)
            {
                if (f(e))
                {
                    // e is a reference here.
                    l.Append(e);
                }
            }
            return l;
        }
    }
}