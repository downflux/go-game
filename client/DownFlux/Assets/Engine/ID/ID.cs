namespace DF.Game.ID
{
    public class Tick : System.IComparable
    {
        private double _tick;
        public Tick(double t) { Double = t; }
        public double Double
        {
            get => _tick;
            private set { _tick = value; }
        }

        public int CompareTo(object other)
        {
            var v = other as Tick;
            return Double.CompareTo(v.Double);
        }

        public static bool operator <(Tick a, Tick b) { return a.Double < b.Double; }
        public static bool operator >(Tick a, Tick b) { return a.Double > b.Double; }

        public static bool operator >=(Tick a, Tick b) { return a.Double >= b.Double; }
        public static bool operator <=(Tick a, Tick b) { return a.Double <= b.Double; }
    }

    public class ID
    {
        private string _id;
        public string String
        {
            get => _id;
            private set { _id = value; }
        }

        public ID(string id) { String = id; }
    }

    public class EntityID : ID
    {
        public EntityID(string id) : base(id) { }
    }
}