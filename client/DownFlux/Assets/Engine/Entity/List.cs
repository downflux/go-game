namespace DF.Game.Entity
{
    public class List
    {
        private System.TimeSpan _acquireTimeout;
        private System.Collections.Generic.Dictionary<
            DF.Game.ID.EntityID,
            DF.Game.Entity.IEntity> _entities;
        private System.Threading.ReaderWriterLock _l;

        public List(System.TimeSpan acquireTimeout)
        {
            _acquireTimeout = acquireTimeout;
            _l = new System.Threading.ReaderWriterLock();
            _entities = new System.Collections.Generic.Dictionary<DF.Game.ID.EntityID, DF.Game.Entity.IEntity>();
        }

        private void Append(IEntity e)
        {
            _l.AcquireWriterLock(_acquireTimeout);
            try
            {
                _entities[e.ID] = e;
            }
            finally
            {
                _l.ReleaseWriterLock();
            }

        }

        protected System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.Entity.IEntity> Entities
        {
            get
            {
                System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.Entity.IEntity> es;
                _l.AcquireReaderLock(_acquireTimeout);
                try
                {
                    es = new System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.Entity.IEntity>(
                        new System.Collections.Generic.List<DF.Game.Entity.IEntity>(_entities.Values));
                }
                finally
                {
                    _l.ReleaseReaderLock();
                }
                return es;
            }
        }

        public void Merge(List o)
        {
            _l.AcquireWriterLock(_acquireTimeout);
            try
            {
                foreach (var e in o.Entities)
                {
                    if (_entities.ContainsKey(e.ID))
                    {

                    }
                    else
                    {
                        Append(e);
                    }
                }
            }
            finally
            {
                _l.ReleaseWriterLock();
            }
        }
    }
}