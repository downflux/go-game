using System.Collections;
using System.Collections.Generic;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.TestTools;

namespace Tests
{
    public class StatusTest
    {
        [Test]
        public void TestStatusTick()
        {
            var s = new DF.Game.Status.Status(DF.Test.Data.StatusTestData.pb);
            Assert.AreEqual(new DF.Game.ID.Tick(10), s.CalculateTick(new System.DateTime(2000, 1, 1, 0, 0, 1)));
        }
    }
}