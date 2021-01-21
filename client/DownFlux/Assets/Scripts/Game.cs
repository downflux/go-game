using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Game : MonoBehaviour
{
    public string server;
    private static System.TimeSpan _serverBootSleepDuration = new System.TimeSpan(System.Convert.ToInt64(1e7));
    private static System.TimeSpan _entityListAcquireDuration = new System.TimeSpan(System.Convert.ToInt64(1e7));
    private DF.Game.Game game;

    // Start is called before the first frame update
    void Start()
    {
        game = new DF.Game.Game(
            server,
            new DF.Game.Config(_serverBootSleepDuration, _entityListAcquireDuration)
        );
    }

    // Update is called once per frame
    void Update()
    {
        
    }
}
