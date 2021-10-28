namespace DF.Game.Utilities.FSM {
    // TODO(minkezhang): Change to Dictionary<State, List<State>> instead.
    using TransitionLookup = System.Collections.ObjectModel.ReadOnlyDictionary<State, State>;

    public class State {
        private int _s;
        public State(int s) {
            _s = s;
        }
    }
    
    public interface IFSM {
        bool To(State a, State b);
        State State();
    }

    public abstract class Base : IFSM {
        private TransitionLookup _transitions;
        private State _state;  // Default state must be UNKNOWN.

        public Base(TransitionLookup transitions, State s) {
            _transitions = transitions;
            _state = s;
        }

        public bool To(State a, State b) {
            if (_transitions.ContainsKey(a) && _transitions[a] == b) {
                _state = b;
                return true;
            }
            return false;
        }

        public abstract State State();
    }
}