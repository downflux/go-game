namespace DF.Game.Exception
{
    public class MergeException : System.ArgumentException
    {
        public MergeException() { }

        public MergeException(string message) : base(message) { }

        public MergeException(string message, System.Exception inner) : base(message, inner) { }
    }

    public class DuplicateKeyException : System.ArgumentException
    {
        public DuplicateKeyException() { }

        public DuplicateKeyException(string message) : base(message) { }

        public DuplicateKeyException(string message, System.Exception inner) : base(message, inner) { }
    }
}