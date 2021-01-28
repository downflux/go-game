using System.Collections;
using System.Collections.Generic;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.TestTools;

namespace Tests
{
    public class DataTest
    {
        [Test]
        public void TestNewData()
        {
            var d = new DF.Game.Data.Data<int>(new List<DF.Game.Data.Datum<int>>());

            Assert.AreEqual(0, d.Count);
        }

        [Test]
        public void TestNewDataWithData()
        {
            int t = 1;
            int v = 10;

            var datum = new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(t), v);
            var d = new DF.Game.Data.Data<int>(new List<DF.Game.Data.Datum<int>> { datum });

            Assert.AreEqual(1, d.Count);
            Assert.AreEqual(
                new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(t), v),
                d[0],
                System.String.Format("t == {0}, v == {1}", d[0].Tick.Double, d[0].Value));
        }

        [Test]
        public void TestBinarySearch()
        {
            var dataList = new List<DF.Game.Data.Datum<int>> {
                new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(1), 1),
                new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(10), 2),
                new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(100), 3)
            };

            var testConfigs = new System.Collections.Generic.List<dynamic>{
                new {
                    name = "SearchPre",
                    t = new DF.Game.ID.Tick(0),
                    want = -1,
                },
                new {
                    name = "SearchAt",
                    t = new DF.Game.ID.Tick(1),
                    want = 0,
                },
                new {
                    name = "SearchBetween",
                    t = new DF.Game.ID.Tick(2),
                    want = -2,
                }
            };

            foreach (var c in testConfigs)
            {
                var d = new DF.Game.Data.Data<int>(dataList);
                Assert.AreEqual(c.want, d.BinarySearch(c.t));
            }
        }

        [Test]
        public void TestDataTruncate()
        {
            var refDataList = new List<DF.Game.Data.Datum<int>> {
                new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(1), 1),
                new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(10), 2),
                new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(100), 3)
            };
            var testConfigs = new System.Collections.Generic.List<dynamic>{
                new {
                    name = "TruncateAll",
                    t = new DF.Game.ID.Tick(0),
                    want = new System.Collections.Generic.List<DF.Game.Data.Datum<int>>{ }
                },
                new {
                    name = "TruncateAllAtTick",
                    t = new DF.Game.ID.Tick(1),
                    want = new System.Collections.Generic.List<DF.Game.Data.Datum<int>>{ }
                },
                new {
                    name = "TruncatePartial",
                    t = new DF.Game.ID.Tick(2),
                    want = new System.Collections.Generic.List<DF.Game.Data.Datum<int>>{
                        refDataList[0]
                    }
                },
                new {
                    name = "TruncateNone",
                    t = new DF.Game.ID.Tick(101),
                    want = new System.Collections.Generic.List<DF.Game.Data.Datum<int>>{
                        refDataList[0],
                        refDataList[1],
                        refDataList[2]
                    }
                }
            };

            foreach (var c in testConfigs)
            {
                var dataList = new List<DF.Game.Data.Datum<int>> {
                    new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(1), 1),
                    new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(10), 2),
                    new DF.Game.Data.Datum<int>(new DF.Game.ID.Tick(100), 3)
                };

                var d = new DF.Game.Data.Data<int>(dataList);
                d.Truncate(c.t);
                Assert.AreEqual(c.want.Count, d.Count);
                for (var i = 0; i < d.Count; i++)
                {
                    Assert.AreEqual(c.want[i], d[i]);
                }
            }
        }
    }
}