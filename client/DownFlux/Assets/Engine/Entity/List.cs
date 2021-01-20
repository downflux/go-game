namespace DF.Game.Entity
{
    public class List
    {
        private System.Collections.Generic.Dictionary<
            DF.Game.ID.EntityID,
            DF.Game.Entity.IEntity> _entities;
        
        public List()
        {
            _entities = new System.Collections.Generic.Dictionary<DF.Game.ID.EntityID, DF.Game.Entity.IEntity>();
        }

        private void Append(IEntity e)
        {
            _entities[e.ID] = e;
        }
        
        protected System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.Entity.IEntity> Entities {
            get => new System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.Entity.IEntity>(
                new System.Collections.Generic.List<DF.Game.Entity.IEntity>(_entities.Values));
        }

        public void Merge(List o) {
            foreach (var e in o.Entities) {
                if (_entities.ContainsKey(e.ID)) {

                } else {
                    Append(e);
                }
            }
        }
    }
}