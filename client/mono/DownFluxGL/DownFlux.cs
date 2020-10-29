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
        private System.Collections.Generic.Dictionary<
          (string, DF.Game.API.Constants.CurveCategory),
          OneOf.OneOf<DF.Curve.LinearMove>> _curves;

        private System.TimeSpan _serverTickDuration;

        public DownFlux(string server)
        {
            _server = server;
            _graphics = new GraphicsDeviceManager(this);
            _curves = new System.Collections.Generic.Dictionary<
              (string, DF.Game.API.Constants.CurveCategory),
              OneOf.OneOf<DF.Curve.LinearMove>>();

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
            }

            // TODO(minkezhang): Make this async task actually have a Wait()
            // somewhere.
            // _c.StreamCurvesLoop(_tid).Start();
            System.Threading.Tasks.Task.Run(() => _c.StreamCurvesLoop(tick));

            // TODO(minkezhang): gGenerate entity name dynamically instead.
            // TODO(minkezhang): Call Move() only when user clicks on map.
            _c.Move(
              tick,
              new System.Collections.Generic.List<string>(){ "example-entity" },
              new DF.Game.API.Data.Position{
                X = 5,
                Y = 5
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
                  System.Console.Error.WriteLine("MATCHED LINEAR MOVE");
                  if (!_curves.ContainsKey((linearMove.EntityID, DF.Game.API.Constants.CurveCategory.Move))) {
                    System.Console.Error.WriteLine("Adding curve for entity: ", linearMove.EntityID);
                    System.Console.Error.WriteLine(linearMove);
                    _curves[(linearMove.EntityID, DF.Game.API.Constants.CurveCategory.Move)] = linearMove;
                  } else {
                    _curves[(linearMove.EntityID, DF.Game.API.Constants.CurveCategory.Move)].AsT0.ReplaceTail(linearMove);
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
                  System.Console.Error.WriteLine(tick);
                  System.Console.Error.WriteLine(p);
                  System.Console.Error.WriteLine(linearMove.Get(30));

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
