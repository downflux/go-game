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