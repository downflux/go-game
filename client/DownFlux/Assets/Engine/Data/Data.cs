[assembly: System.Runtime.CompilerServices.InternalsVisibleTo("EngineTestAssembly")]

namespace DF.Game.Data
{
    public class Datum<T> : System.IComparable
    {
        private DF.Game.ID.Tick _t;
        private T _v;

        public Datum(DF.Game.ID.Tick tick, T value) : this(tick)
        {
            Value = value;
        }

        public Datum(DF.Game.ID.Tick tick)
        {
            Tick = tick;
        }

        public DF.Game.ID.Tick Tick
        {
            get => _t;
            private set { _t = value; }
        }

        public T Value
        {
            get => _v;
            private set { _v = value; }
        }

        public override bool Equals(object obj)
        {
            var o = obj as Datum<T>;
            return Tick.Equals(o.Tick) && Value.Equals(o.Value);
        }

        public override int GetHashCode()
        {
            int hash = 13;
            hash = (hash * 7) + Tick.GetHashCode();
            hash = (hash * 7) + Value.GetHashCode();
            return hash;
        }

        public int CompareTo(object other) { return _t.CompareTo((other as Datum<T>)._t); }
    }

    public class Data<T>
    {
        private System.Collections.Generic.List<Datum<T>> _data;

        public Data(System.Collections.Generic.List<Datum<T>> data)
        {
            // We are assuming input data is already sorted by ticks. Input
            // data is consumed directly from a GameState proto.
            _data = data;
        }

        public int Count => _data.Count;

        public Datum<T> this[int i]
        {
            get => _data[i];
        }

        public int BinarySearch(DF.Game.ID.Tick tick)
        {
            return _data.BinarySearch(new Datum<T>(tick));
        }

        // Truncate deletes all data in this struct after the input tick.
        internal void Truncate(DF.Game.ID.Tick tick)
        {
            var i = BinarySearch(tick);

            if (i == ~_data.Count)
            {
                return;
            }

            if (i < 0)
            {
                i = ~i;
            }

            _data.RemoveRange(i, _data.Count - i);
        }

        // Merge overwrites an instance's own data with the data from the
        // target. We assume data in the target is already in order.
        public void Merge(Data<T> other)
        {
            if (other == null || other.Count == 0)
            {
                return;
            }

            Truncate(other[0].Tick);

            _data.AddRange(other._data);
        }
    }
}