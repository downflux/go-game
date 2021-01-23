using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Game : MonoBehaviour
{
    public string server;
    private static System.TimeSpan _serverBootSleepDuration = new System.TimeSpan(System.Convert.ToInt64(1e7));
    private static System.TimeSpan _entityListAcquireDuration = new System.TimeSpan(System.Convert.ToInt64(1e7));
    private System.Threading.CancellationTokenSource _cancellation;
    private DF.Game.Game _g;

    void Start()
    {
        var config = new DF.Game.Config(
            serverBootSleepDuration: _serverBootSleepDuration,
            entityListAcquireDuration: _entityListAcquireDuration
        );
        _cancellation = new System.Threading.CancellationTokenSource();
        _g = new DF.Game.Game(server, config, _cancellation.Token);
    }

    void Update()
    {
    }

    // OnApplicationQuit is called when the game stops. This may not be called
    // on different platforms. See
    // https://docs.unity3d.com/ScriptReference/MonoBehaviour.OnApplicationQuit.html.
    void OnApplicationQuit()
    {
        _cancellation.Cancel();
    }
}