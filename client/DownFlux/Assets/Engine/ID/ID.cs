namespace DF
{
    namespace Game
    {
        class ID
        {
            private string _id;
            public string String
            {
                get => _id;
                private set { _id = value; }
            }
            public ID(string id) { _id = id; }
        }

        class EntityID : ID
        {
            public EntityID(string id) : base(id) { }
        }
    }
}