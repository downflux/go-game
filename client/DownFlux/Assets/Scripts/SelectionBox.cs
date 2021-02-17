using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class SelectionBox : MonoBehaviour
{
    public RectTransform selection;
    private Vector2 start;
    private bool isDown;
    private const int primaryButton = 0;

    // Start is called before the first frame update
    void Start()
    {
        isDown = false;
    }

    // Update is called once per frame
    void Update()
    {
        if (Input.GetMouseButtonDown(primaryButton))
        {
            isDown = true;
            start = Input.mousePosition;
            selection.position = Input.mousePosition;
        }
        else if (Input.GetMouseButtonUp(primaryButton))
        {
            isDown = false;
            selection.gameObject.SetActive(false);
            var p0 = selection.anchoredPosition - selection.sizeDelta / 2;  // lower bound
            var p1 = selection.anchoredPosition + selection.sizeDelta / 2;  // upper bound
        }
        else if (isDown)
        {
            var mousePos = Input.mousePosition;
            var x = mousePos.x - start.x;
            var y = mousePos.y - start.y;
            selection.sizeDelta = new Vector2(Mathf.Abs(x), Mathf.Abs(y));
            selection.anchoredPosition = start + new Vector2(x / 2, y / 2);
        }
    }
}
