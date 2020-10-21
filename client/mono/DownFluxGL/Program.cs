using System;

namespace DownFluxGL
{
    public static class Program
    {
        [STAThread]
        static void Main()
        {
            using (var game = new DownFlux())
                game.Run();
        }
    }
}
