using System.Collections;
using System.Collections.Generic;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.TestTools;

namespace Tests
{
    public class LinearPositionTest
    {
        [Test]
        public void TestMerge()
        {
            var curve = new DF.Game.Curve.LinearPosition.LinearPosition(DF.Test.Data.LinearPositionTestData.pb);
            curve.Merge(
                new DF.Game.Curve.LinearPosition.LinearPosition(DF.Test.Data.LinearPositionTestData.mergePB));

            var testConfigs = new System.Collections.Generic.List<dynamic>{
                new {
                    name = "Before",
                    t = new DF.Game.ID.Tick(-1),
                    want = new DF.Game.API.Data.Position{X = 0, Y = 0}
                },
                new {
                    name = "BeforeNewCurveInterpolation",
                    t = new DF.Game.ID.Tick(.1),
                    want = new DF.Game.API.Data.Position{X = .1, Y = .1}
                },
                new {
                    name = "NewCurveBegin",
                    t = new DF.Game.ID.Tick(1),
                    want = new DF.Game.API.Data.Position{X = 1, Y = 1}
                },
                new {
                    name = "NewCurveBetween",
                    t = new DF.Game.ID.Tick(2),
                    want = new DF.Game.API.Data.Position{X = 1.1, Y = 1.1}
                },
                new {
                    name = "After",
                    t = new DF.Game.ID.Tick(11),
                    want = new DF.Game.API.Data.Position{X = 2, Y = 2}
                },
            };

            foreach (var c in testConfigs)
            {
                Assert.AreEqual(c.want, curve.Get(c.t));
            }
        }

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
