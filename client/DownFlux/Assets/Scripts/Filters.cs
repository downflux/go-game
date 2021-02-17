using System.Collections;
using System.Collections.Generic;
using UnityEngine;
namespace DF.Unity
{
    public static class Filters
    {
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