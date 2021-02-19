using System.Collections;
using System.Collections.Generic;
using UnityEngine;

namespace DF.Utilities.IntentState
{
    enum _StateLookup
    {
        Unknown = 0
    }

    public class IntentState : DF.Game.Utilities.FSM.Base
    {
        private static System.Collections.ObjectModel.ReadOnlyDictionary<
            _StateLookup, DF.Game.Utilities.FSM.State> _lookup = new System.Collections.ObjectModel.ReadOnlyDictionary<
                _StateLookup, DF.Game.Utilities.FSM.State>(
                new Dictionary<_StateLookup, DF.Game.Utilities.FSM.State>()
                {
                    { _StateLookup.Unknown, new DF.Game.Utilities.FSM.State((int) _StateLookup.Unknown) }
                }
            );

        private static System.Collections.ObjectModel.ReadOnlyDictionary<
            DF.Game.Utilities.FSM.State, DF.Game.Utilities.FSM.State> _transitions = new System.Collections.ObjectModel.ReadOnlyDictionary<
                DF.Game.Utilities.FSM.State, DF.Game.Utilities.FSM.State>(
            new Dictionary<DF.Game.Utilities.FSM.State, DF.Game.Utilities.FSM.State>()
            {
            }
        );

        public IntentState(DF.Game.Utilities.FSM.State s) : base(transitions: _transitions, s: s) { }
        public override DF.Game.Utilities.FSM.State State()
        {
            if (Input.GetMouseButtonDown(0)) {
                return _lookup[_StateLookup.Unknown];
            }
            return _lookup[_StateLookup.Unknown];
        }
    }
}