namespace DF.Game.Entity
{
    public interface IEntity
    {
        DF.Game.API.Constants.EntityType Type { get; }
        DF.Game.ID.EntityID ID { get; }
        DF.Game.Curve.List Curves { get; }
    }
}