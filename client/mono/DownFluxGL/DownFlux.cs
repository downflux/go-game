using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;
using Microsoft.Xna.Framework.Input;

namespace DownFluxGL
{
    public class DownFlux : Game
    {
        private static int tileWidth = 100;

        Texture2D ballTexture;

        private GraphicsDeviceManager _graphics;
        private SpriteBatch _spriteBatch;

        // DownFlux server API
        private DF.Client.Client _c;
        private string _server;
        private string _tid;
        private System.Collections.Generic.Dictionary<
          string, OneOf.OneOf<DF.Curve.LinearMove>> _curves;

        // TODO(minkezhang): Get this on Initialize() from server.
        private System.TimeSpan _serverTickDuration;

        public DownFlux(string server, string tickID)
        {
            _server = server;
            _tid = tickID;
            _graphics = new GraphicsDeviceManager(this);
            _curves = new System.Collections.Generic.Dictionary<
              string, OneOf.OneOf<DF.Curve.LinearMove>>();

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
            _c.Connect(_tid);
            _serverTickDuration = new System.TimeSpan((long) 1e6); // 100 ms

            // TODO(minkezhang): Make this async task actually have a Wait()
            // somewhere.
            _c.StreamCurvesLoop(_tid).Start();

            // TODO(minkezhang): Remove this.
            _c.Move(
              "",
              new System.Collections.Generic.List<string>(){ "example-entity" },
              new DF.Game.API.Data.Position{
                X = 5 * tileWidth,
                Y = 5 * tileWidth
              },
              DF.Game.API.Constants.MoveType.Forward
            );

            base.Initialize();
        }

        protected override void LoadContent()
        {
            _spriteBatch = new SpriteBatch(GraphicsDevice);

            ballTexture = Content.Load<Texture2D>("Assets/ball");
        }

        protected override void Update(GameTime gameTime)
        {
            foreach (var d in _c.Data) {
              d.Item2.Switch(
                linearMove => {
                  if (!_curves.ContainsKey(linearMove.ID)) {
                    _curves[linearMove.ID] = linearMove;
                  } else {
                    _curves[linearMove.ID].AsT0.ReplaceTail(linearMove);
                  }
                }
              );
            }

            base.Update(gameTime);
        }

        protected override void Draw(GameTime gameTime)
        {
            var tick = gameTime.TotalGameTime / _serverTickDuration;

            GraphicsDevice.Clear(Color.CornflowerBlue);

            _spriteBatch.Begin();

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
