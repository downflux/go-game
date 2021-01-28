using System.Collections;
using System.Collections.Generic;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.TestTools;

namespace Tests
{
    public sealed class MockCurve : DF.Game.Curve.Base
    {
        public MockCurve(DF.Game.API.Data.Curve pb) : base(pb) { }

    }

    public class CurveTest
    {
        [Test]
        public void TestNewCurve()
        {
            var eid = "entity-id";
            var tick = 10;
            var type = DF.Game.API.Constants.CurveType.LinearMove;
            var property = DF.Game.API.Constants.EntityProperty.Position;

            var c = new MockCurve(new DF.Game.API.Data.Curve
            {
                EntityId = eid,
                Tick = tick,
                Type = type,
                Property = property
            });

            Assert.AreEqual(eid, c.ID.String);
            Assert.AreEqual(tick, c.Tick.Double);
            Assert.AreEqual(type, c.Type);
            Assert.AreEqual(property, c.Property);
        }
    }
}
