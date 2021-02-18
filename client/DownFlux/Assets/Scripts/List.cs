using System.Collections;
using System.Collections.Generic;
using UnityEngine;
namespace DF.Unity
{
    public delegate bool F(DF.Unity.List.Entity e);

    public class List : MonoBehaviour
    {

        public GameObject TankModel;

        public class Entity
        {
            private GameObject _o;
            private DF.Game.Entity.Entity _e;

            public Entity(GameObject o, DF.Game.Entity.Entity e)
            {
                O = o;
                E = e;
            }

            public GameObject O
            {
                get => _o;
                private set { _o = value; }
            }

            public DF.Game.Entity.Entity E
            {
                get => _e;
                private set { _e = value; }
            }
        }

        private System.Collections.Generic.Dictionary<DF.Game.ID.EntityID, Entity> _entities;

        public Entity Get(DF.Game.ID.EntityID id) => _entities[id];

        public void Append(DF.Game.Entity.Entity entity)
        {
            switch (entity.Type)
            {
                case DF.Game.API.Constants.EntityType.Tank:
                    if (_entities.ContainsKey(entity.ID))
                    {
                        return;
                    }
                    _entities[entity.ID] = new Entity(
                        Instantiate(TankModel, transform.position, transform.rotation),
                        entity
                    );
                    break;
            }

        }

        // Start is called before the first frame update
        void Start()
        {
            _entities = new System.Collections.Generic.Dictionary<DF.Game.ID.EntityID, Entity>();
        }

        // Update is called once per frame
        void Update()
        {
            foreach (var e in _entities.Values)
            {
                var (ok, c) = e.E.Curves.Curve(DF.Game.API.Constants.EntityProperty.Position).TryGetLinearPosition();
                if (ok)
                {
                    var p = c.Get(GetComponent<DF.Unity.Game>().Tick);
                    e.O.transform.position = new Vector3((float)p.X, 0, (float)p.Y);
                }
            }
        }

        public List<DF.Game.ID.EntityID> Filter(F f)
        {
            var l = new List<DF.Game.ID.EntityID>();
            foreach (var e in _entities.Values)
            {
                if (f(e))
                {
                    // e is a reference here.
                    l.Add(e.E.ID);
                }
            }
            return l;
        }
    }
}