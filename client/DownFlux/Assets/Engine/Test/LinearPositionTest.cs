using System.Collections;
using System.Collections.Generic;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.TestTools;

namespace Tests
{
    public class LinearPositionTest
    {
        // TODO(minkezhang): Add MergeTest.

        [Test]
        public void TestGet()
        {
            var curve = new DF.Game.Curve.LinearPosition.LinearPosition(DF.Test.Data.LinearPositionTestData.pb);

            var testConfigs = new System.Collections.Generic.List<dynamic>{
                new {
                    name = "Before",
                    t = new DF.Game.ID.Tick(-1),
                    want = new DF.Game.API.Data.Position{X = 0, Y = 0}
                },
                new {
                    name = "Between",
                    t = new DF.Game.ID.Tick(1),
                    want = new DF.Game.API.Data.Position{X = 0.1, Y = 0.1}
                },
                new {
                    name = "At",
                    t = new DF.Game.ID.Tick(10),
                    want = new DF.Game.API.Data.Position{X = 1, Y = 1}
                },
                new {
                    name = "After",
                    t = new DF.Game.ID.Tick(11),
                    want = new DF.Game.API.Data.Position{X = 1, Y = 1}
                },
            };

            foreach (var c in testConfigs)
            {
                Assert.AreEqual(c.want, curve.Get(c.t));
            }
        }
    }
}
