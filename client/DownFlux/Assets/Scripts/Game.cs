using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Game : MonoBehaviour
{
    public string server;
    private static System.TimeSpan _serverBootSleepDuration = new System.TimeSpan(System.Convert.ToInt64(1e7));
    private static DF.Game.ID.Tick _updateTickDelay = new DF.Game.ID.Tick(0.5);
    private System.Threading.CancellationTokenSource _cancellation;
    private DF.Game.Game _g;
    private DF.Game.ID.Tick _lastMergeTick;

    private void Quit()
    {
        // save any game data here
#if UNITY_EDITOR
        // Application.Quit() does not work in the editor so
        // UnityEditor.EditorApplication.isPlaying need to be set to false to end the game
        UnityEditor.EditorApplication.isPlaying = false;
#else
        Application.Quit();
#endif
    }
    void Start()
    {
        var config = new DF.Game.Config(
            serverBootSleepDuration: _serverBootSleepDuration,
            listenerAcquireTimeout: _serverBootSleepDuration,
            updateTickDelay: _updateTickDelay
        );
        _cancellation = new System.Threading.CancellationTokenSource();

        try
        {
            _g = new DF.Game.Game(server, config, _cancellation.Token);
        }
        catch (Grpc.Core.RpcException)
        {
            Quit();
        }
        _lastMergeTick = new DF.Game.ID.Tick(_g.Tick);
    }

    void Update()
    {
        var t = new DF.Game.ID.Tick(_g.Tick);
        if (t - _lastMergeTick > _updateTickDelay)
        {
            try {
                _g.Merge();
                _lastMergeTick = t;
            } catch (Grpc.Core.RpcException) {
                Quit();
            }
        }
        // TODO(minkezhang): Update client state with position.
    }

    // OnApplicationQuit is called when the game stops. This may not be called
    // on different platforms. See
    // https://docs.unity3d.com/ScriptReference/MonoBehaviour.OnApplicationQuit.html.
    void OnApplicationQuit()
    {
        _cancellation.Cancel();
    }
}