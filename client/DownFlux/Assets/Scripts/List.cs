using System.Collections;
using System.Collections.Generic;
using UnityEngine;
namespace DF.Unity
{
    public delegate bool F(DF.Unity.List.Entity e);

    public class List : MonoBehaviour
    {
        private System.Collections.Generic.Dictionary<
            DF.Game.API.Constants.EntityType, GameObject> _modelLookup;
        public GameObject TankModel;
        public GameObject ShellModel;

        public class Entity
        {
            // TODO(minkezhang): Consider ways to encapsulate entity data in
            // the GameObject itself instead.
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
        private System.Collections.Generic.Dictionary<int, DF.Game.ID.EntityID> _instanceIDs;

        public Entity Get(DF.Game.ID.EntityID id) => _entities[id];

        public void Append(DF.Game.Entity.Entity entity)
        {
            if (
                _modelLookup.ContainsKey(entity.Type) && !(
                    _entities.ContainsKey(entity.ID)))
            {
                var e = new Entity(
                    Instantiate(
                        _modelLookup[entity.Type],
                        transform.position,
                        transform.rotation),
                    entity
                );
                e.O.layer = LayerMask.NameToLayer("Entities");
                _entities[entity.ID] = e;
                _instanceIDs[e.O.GetInstanceID()] = e.E.ID;
            }
        }
        public Entity GetByInstanceID(int id) => GetByEntityID(_instanceIDs[id]);
        public Entity GetByEntityID(DF.Game.ID.EntityID eid) => _entities[eid];


        // Start is called before the first frame update
        void Start()
        {
            _entities = new System.Collections.Generic.Dictionary<DF.Game.ID.EntityID, Entity>();
            _instanceIDs = new System.Collections.Generic.Dictionary<int, DF.Game.ID.EntityID>();

            _modelLookup = new System.Collections.Generic.Dictionary<
            DF.Game.API.Constants.EntityType, GameObject>{
                { DF.Game.API.Constants.EntityType.Tank, TankModel },
                { DF.Game.API.Constants.EntityType.TankProjectile, ShellModel }
            };
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