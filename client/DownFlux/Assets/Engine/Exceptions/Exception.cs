namespace DF.Game.Exception
{
    public class MergeException : System.ArgumentException
    {
        public MergeException() { }

        public MergeException(string message) : base(message) { }

        public MergeException(string message, System.Exception inner) : base(message, inner) { }

    }
}