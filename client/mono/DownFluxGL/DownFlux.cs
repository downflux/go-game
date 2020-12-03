using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;
using Microsoft.Xna.Framework.Input;

namespace DownFluxGL
{
    public class DownFlux : Game
    {
        private DF.Input.Event.Mouse _mouseState;

	private bool _mouseDown;
        private Microsoft.Xna.Framework.Point m0, m1;

        private static int tileWidth = 100;

        Texture2D ballTexture;
        Texture2D boxTexture;
        private Microsoft.Xna.Framework.Rectangle selectionBox;
        private System.Collections.Generic.HashSet<string> _selectedEntities;

        private GraphicsDeviceManager _graphics;
        private SpriteBatch _spriteBatch;

        // DownFlux server API
        private DF.Client.Client _c;
        private string _server;

        // Add-only collections
        private System.Collections.Generic.Dictionary<
          string, DF.Entity.Entity> _entities;
        private System.Collections.Generic.Dictionary<
          (string, DF.Game.API.Constants.CurveCategory),
          DF.Curve.Curve> _curves;

        private System.TimeSpan _serverTickDuration;
        private System.DateTime _serverStartTime;

        public DownFlux(string server)
        {
            _server = server;
            _graphics = new GraphicsDeviceManager(this);
            _curves = new System.Collections.Generic.Dictionary<
              (string, DF.Game.API.Constants.CurveCategory),
              DF.Curve.Curve>();
            _entities = new System.Collections.Generic.Dictionary<
              string, DF.Entity.Entity>();
            _selectedEntities = new System.Collections.Generic.HashSet<string>();

            _mouseState = new DF.Input.Event.Mouse();

            // TODO(minkezhang): Figure out why we can't control actual window
            // size here.
            _graphics.PreferredBackBufferWidth = 10 * tileWidth;
            _graphics.PreferredBackBufferHeight = 10 * tileWidth;
            _graphics.ApplyChanges();

            Content.RootDirectory = "Content";
            IsMouseVisible = true;
        }

        protected override void Initialize()
        {
            _c = new DF.Client.Client(
              new Grpc.Core.Channel(
                _server, Grpc.Core.ChannelCredentials.Insecure));
            _c.Connect();

            bool isStarted = false;
            double tick = 0;
            // TODO(minkezhang): Import map.
            while (!isStarted) {
              var status = _c.GetStatus();
              _serverTickDuration = status.TickDuration.ToTimeSpan();
              isStarted = status.IsStarted;
              tick = status.Tick;
              _serverStartTime = status.StartTime.ToDateTime().ToLocalTime();
            }

            // TODO(minkezhang): Make this async task actually have a Wait()
            // in destructor.
            // _c.StreamDataLoop(_tid).Start();
            System.Threading.Tasks.Task.Run(() => _c.StreamDataLoop(tick));

            base.Initialize();
        }

        protected override void LoadContent()
        {
            _spriteBatch = new SpriteBatch(GraphicsDevice);

            ballTexture = Content.Load<Texture2D>("Assets/ball");
            boxTexture = Content.Load<Texture2D>("Assets/blackpixel");
        }

        protected override void Update(GameTime gameTime)
        {
            var ms = _mouseState.GetState();
            var recalculateBox = false;
            var selectEntities = false;
            var doAction = false;
            _mouseDown = ms.IsHeld();

            if (ms.IsPressed()) {
              System.Console.Error.WriteLine("Button is pressed");
              m0 = new Microsoft.Xna.Framework.Point(ms.State.X, ms.State.Y);
              m1 = new Microsoft.Xna.Framework.Point(ms.State.X, ms.State.Y);
            } else if (ms.IsHeld()) {
              System.Console.Error.WriteLine("Button is held");
              m1 = new Microsoft.Xna.Framework.Point(ms.State.X, ms.State.Y);
            } else if (ms.IsReleased()) {
              System.Console.Error.WriteLine("Button is released");
              // TODO(minkezhang): Add some leeway here -- fast clicks may
              // accidentally turn into a drag motion.
              if (m0 == m1) {
                System.Console.Error.WriteLine("Interpreting as click");
                doAction = true;
              } else {
                System.Console.Error.WriteLine("Interpreting as drag");
                selectEntities = true;
              }
            }
            recalculateBox = ms.IsHeld();

            if (selectEntities) {
              System.Console.Error.WriteLine("clearing entities");
              _selectedEntities.Clear();
            }

            if (recalculateBox) {
              selectionBox = new Microsoft.Xna.Framework.Rectangle(
                System.Math.Min(m0.X, m1.X),
                System.Math.Min(m0.Y, m1.Y),
                System.Math.Abs(m1.X - m0.X),
                System.Math.Abs(m1.Y - m0.Y));
            }

            var tick = (System.DateTime.Now - _serverStartTime) / _serverTickDuration;

            // TODO(minkezhang): Determine if we need d to be a CurveStream
            // object (with associated tick).
            foreach (var d in _c.Data) {
              d.Item2.Switch(
                linearMove => {
                  if (!_curves.ContainsKey((linearMove.EntityID, DF.Game.API.Constants.CurveCategory.Move))) {
                    System.Console.Error.WriteLine(linearMove);
                    _curves[(linearMove.EntityID, DF.Game.API.Constants.CurveCategory.Move)] = linearMove;
                  } else {
                    _curves[(linearMove.EntityID, DF.Game.API.Constants.CurveCategory.Move)].AsT0.ReplaceTail(linearMove);
                  }
                }
              );
            }

            foreach (var e in _c.Entities) {
              e.Item2.Switch(
                simpleEntity => {
                  if (!_entities.ContainsKey(simpleEntity.ID)) {
                    _entities[simpleEntity.ID] = simpleEntity;
                  }
                }
              );
            }

            foreach (var e in _entities) {
              e.Value.Switch(
                simpleEntity => {
                  DF.Curve.Curve c;
                  if (selectEntities && !_selectedEntities.Contains(simpleEntity.ID)) {
                    _curves.TryGetValue((simpleEntity.ID, DF.Game.API.Constants.CurveCategory.Move), out c);
                    c.Switch(
                      linearMove => {
                        var p = linearMove.Get(tick);
                        System.Console.Error.WriteLine(p);
                        System.Console.Error.WriteLine(selectionBox);
                        if (selectionBox.Contains(new Microsoft.Xna.Framework.Vector2((float) p.X * tileWidth, (float) p.Y * tileWidth))) {
                          _selectedEntities.Add(simpleEntity.ID);
                        }
                      }
                    );
                  };
                }
              );
            }

            if (doAction) {
              System.Console.Error.WriteLine("ENTITIES: ", _selectedEntities);
              if (_selectedEntities.Count > 0) {
                _c.Move(
                  tick,
                    new System.Collections.Generic.List<string>(_selectedEntities),
                    new DF.Game.API.Data.Position{
                      X = ms.State.X / tileWidth,
                      Y = ms.State.Y / tileWidth
                    },
                    DF.Game.API.Constants.MoveType.Forward
                );
              }
            }
            base.Update(gameTime);
        }

        protected override void Draw(GameTime gameTime)
        {
            var tick = (System.DateTime.Now - _serverStartTime) / _serverTickDuration;

            GraphicsDevice.Clear(Color.CornflowerBlue);

            _spriteBatch.Begin();

            if (_mouseDown) {
              _spriteBatch.Draw(boxTexture, selectionBox, Color.White);
            }

            foreach (var c in _curves) {
              c.Value.Switch(
                linearMove => {
                  var p = linearMove.Get(tick);

                  _spriteBatch.Draw(
                    ballTexture,
                    new Vector2((float) p.X * tileWidth, (float) p.Y * tileWidth),
                    null,
                    Color.White,
                    0f,
                    new Vector2(ballTexture.Width / 2, ballTexture.Height / 2),
                    Vector2.One,
                    SpriteEffects.None,
                    0f
                  );
                }
              );
            }

            _spriteBatch.End();

            base.Draw(gameTime);
        }
    }
}
