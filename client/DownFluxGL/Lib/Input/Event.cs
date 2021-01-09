namespace DF {
  namespace Input {
    namespace Event {
      public struct Mouse {
        private Microsoft.Xna.Framework.Input.MouseState _cachedState;
        public MouseState GetState() {
          var state = Microsoft.Xna.Framework.Input.Mouse.GetState();
          var newState = new MouseState(state, _cachedState);
          _cachedState = state;
          return newState;
        }
      }
      public struct MouseState {
        private Microsoft.Xna.Framework.Input.MouseState _state, _cachedState;
        public MouseState(
          Microsoft.Xna.Framework.Input.MouseState state,
          Microsoft.Xna.Framework.Input.MouseState cachedState) {
            _state = state;
            _cachedState = cachedState;
        }

        public Microsoft.Xna.Framework.Input.MouseState State {
          get => _state;
        }

        public bool IsReleased() {
          return (
            _cachedState.LeftButton == Microsoft.Xna.Framework.Input.ButtonState.Pressed) && (
            _state.LeftButton == Microsoft.Xna.Framework.Input.ButtonState.Released);
        }
        public bool IsPressed() {
          return (
            _cachedState.LeftButton == Microsoft.Xna.Framework.Input.ButtonState.Released) && (
            _state.LeftButton == Microsoft.Xna.Framework.Input.ButtonState.Pressed);
        }
        public bool IsHeld() {
          return (
            _cachedState.LeftButton == Microsoft.Xna.Framework.Input.ButtonState.Pressed) && (
            _state.LeftButton == Microsoft.Xna.Framework.Input.ButtonState.Pressed);
        }
      }
    }
  }
}
