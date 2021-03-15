using System.Collections;
using System.Collections.Generic;
using UnityEngine;
namespace DF.Unity
{
    public static class Filters
    {
        public static DF.Unity.F And(params DF.Unity.F[] fs)
        {
            return new DF.Unity.F((DF.Unity.List.Entity e) =>
            {
                foreach (var f in fs)
                {
                    // Operations may be expensive -- we want to break early to save on cycles.
                    if (!f(e)) { return false; }
                }
                return true;
            });
        }
        public static DF.Unity.F Not(DF.Unity.F a)
        {
            return new DF.Unity.F((DF.Unity.List.Entity e) => !a(e));
        }
        public static DF.Unity.F FilterByEntityTypes(params DF.Game.API.Constants.EntityType[] ets)
        {
            return new DF.Unity.F((DF.Unity.List.Entity e) =>
            {
                var t = e.E.Type;
                foreach (var et in ets)
                {
                    if (et == t)
                    {
                        return true;
                    }
                }
                return false;
            });
        }
        public static DF.Unity.F FilterByProjectedPosition(Vector2 p0, Vector2 p1, Camera camera)
        {
            return new DF.Unity.F((DF.Unity.List.Entity e) =>
            {
                // TODO(minkezhang): Implement API.Constants.Implements instead, e.g. Moveable, etc.
                if (!e.E.Curves.Properties.Contains(DF.Game.API.Constants.EntityProperty.Position))
                {
                    return false;
                }
                var p = camera.WorldToScreenPoint(e.O.transform.position);
                return (p0.x <= p.x && p.x <= p1.x && p0.y <= p.y && p.y <= p1.y);
            });
        }
    }
}