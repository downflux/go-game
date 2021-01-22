using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Game : MonoBehaviour
{
    public string server;
    private static System.TimeSpan _serverBootSleepDuration = new System.TimeSpan(System.Convert.ToInt64(1e7));
    private static System.TimeSpan _entityListAcquireDuration = new System.TimeSpan(System.Convert.ToInt64(1e7));
    private System.Threading.CancellationTokenSource _cancellation;
    private DF.Game.Game game;

    void Start()
    {
        var config = new DF.Game.Config(
            serverBootSleepDuration: _serverBootSleepDuration,
            entityListAcquireDuration: _entityListAcquireDuration
        );
        _cancellation = new System.Threading.CancellationTokenSource();
        game = new DF.Game.Game(server, config, _cancellation.Token);
    }

    void Update()
    {

    }
}
