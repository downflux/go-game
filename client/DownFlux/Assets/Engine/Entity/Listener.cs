namespace DF.Game.Entity.Listener {
    // LastState overwrites the internal cache with new gRPC server updates.
    //
    // For (playable) game clients, we don't need to rewind into the past, so
    // there's no point in storing corrected past data.
    public class LastState {
        private System.Threading.ReaderWriterLock _l;
        private DF.Game.ID.Tick _serverTick;
        private DF.Game.Entity.List _buffer;
        private System.TimeSpan _acquireTimeout;

        public LastState(System.TimeSpan acquireTimeout) {
            _l = new System.Threading.ReaderWriterLock();
            _acquireTimeout = acquireTimeout;
            _serverTick = new DF.Game.ID.Tick(0);
        }

        public DF.Game.Entity.List Pop() {
            DF.Game.Entity.List v;
            try {
                _l.AcquireReaderLock(_acquireTimeout);
                v = _buffer;
            } finally {
                _l.ReleaseReaderLock();
            }
            return v;
        }

        public void Set(DF.Game.ID.Tick t, DF.Game.Entity.List v) {
            try {
                _l.AcquireWriterLock(_acquireTimeout);
                if (t >= _serverTick) {
                    _serverTick = t;
                    _buffer = v;
                }
            } finally {
                _l.ReleaseWriterLock();
            }
        }

        public void Listen(DF.Game.API.API.StreamDataResponse pb) {
            Set(new DF.Game.ID.Tick(pb.Tick), new DF.Game.Entity.List(pb));
        }
    }
}